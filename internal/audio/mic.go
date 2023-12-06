package audio

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/signal"
	"time"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/moutend/go-wav"
	"github.com/moutend/go-wca/pkg/wca"
)

func Capture(ch chan []byte) {
	err := run(ch)
	if err != nil {
		panic(err)
	}
}

type DurationFlag struct {
	Value time.Duration
}

type FilenameFlag struct {
	Value string
}

func run(ch chan []byte) (err error) {
	for {
		var durationFlag DurationFlag
		var file []byte
		var audio *wav.File

		durationFlag.Value = 500 * time.Millisecond

		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt)
		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			select {
			case <-signalChan:
				cancel()
			}
			return
		}()

		if audio, err = captureSharedTimerDriven(ctx, durationFlag.Value); err != nil {
			continue
		}
		if file, err = wav.Marshal(audio); err != nil {
			return
		}

		if ch != nil {
			ch <- file
		}
	}
}

func captureSharedTimerDriven(ctx context.Context, duration time.Duration) (audio *wav.File, err error) {
	if err = ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED); err != nil {
		return
	}
	defer ole.CoUninitialize()

	var mmde *wca.IMMDeviceEnumerator
	if err = wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &mmde); err != nil {
		return
	}
	defer mmde.Release()

	var mmd *wca.IMMDevice
	if err = mmde.GetDefaultAudioEndpoint(wca.ECapture, wca.EConsole, &mmd); err != nil {
		return
	}
	defer mmd.Release()

	var ps *wca.IPropertyStore
	if err = mmd.OpenPropertyStore(wca.STGM_READ, &ps); err != nil {
		return
	}
	defer ps.Release()

	var pv wca.PROPVARIANT
	if err = ps.GetValue(&wca.PKEY_Device_FriendlyName, &pv); err != nil {
		return
	}

	var ac *wca.IAudioClient
	if err = mmd.Activate(wca.IID_IAudioClient, wca.CLSCTX_ALL, nil, &ac); err != nil {
		return
	}
	defer ac.Release()

	var wfx *wca.WAVEFORMATEX
	if err = ac.GetMixFormat(&wfx); err != nil {
		return
	}
	defer ole.CoTaskMemFree(uintptr(unsafe.Pointer(wfx)))

	wfx.WFormatTag = 1
	wfx.NBlockAlign = (wfx.WBitsPerSample / 8) * wfx.NChannels
	wfx.NAvgBytesPerSec = wfx.NSamplesPerSec * uint32(wfx.NBlockAlign)
	wfx.CbSize = 0

	if audio, err = wav.New(int(wfx.NSamplesPerSec), int(wfx.WBitsPerSample), int(wfx.NChannels)); err != nil {
		return
	}

	var defaultPeriod wca.REFERENCE_TIME
	var minimumPeriod wca.REFERENCE_TIME
	var latency time.Duration
	if err = ac.GetDevicePeriod(&defaultPeriod, &minimumPeriod); err != nil {
		return
	}
	latency = time.Duration(int(defaultPeriod) * 100)

	if err = ac.Initialize(wca.AUDCLNT_SHAREMODE_SHARED, 0, defaultPeriod, 0, wfx, nil); err != nil {
		return
	}

	var bufferFrameSize uint32
	if err = ac.GetBufferSize(&bufferFrameSize); err != nil {
		return
	}

	var acc *wca.IAudioCaptureClient
	if err = ac.GetService(wca.IID_IAudioCaptureClient, &acc); err != nil {
		return
	}
	defer acc.Release()

	if err = ac.Start(); err != nil {
		return
	}

	var output = []byte{}
	var offset int
	var isCapturing bool = true
	var currentDuration time.Duration
	var b *byte
	var data *byte
	var availableFrameSize uint32
	var flags uint32
	var devicePosition uint64
	var qcpPosition uint64

	time.Sleep(latency)

	for {
		if !isCapturing {
			break
		}
		select {
		case <-ctx.Done():
			isCapturing = false
			break
		default:
			// Wait for buffering.
			time.Sleep(latency / 2)

			currentDuration = time.Duration(float64(offset) / float64(wfx.WBitsPerSample/8) / float64(wfx.NChannels) / float64(wfx.NSamplesPerSec) * float64(time.Second))
			if duration != 0 && currentDuration > duration {
				isCapturing = false
				break
			}
			if err = acc.GetBuffer(&data, &availableFrameSize, &flags, &devicePosition, &qcpPosition); err != nil {
				continue
			}
			if availableFrameSize == 0 {
				continue
			}

			start := unsafe.Pointer(data)
			lim := int(availableFrameSize) * int(wfx.NBlockAlign)
			buf := make([]byte, lim)

			for n := 0; n < lim; n++ {
				b = (*byte)(unsafe.Pointer(uintptr(start) + uintptr(n)))
				buf[n] = *b
			}
			offset += lim
			output = append(output, buf...)

			if err = acc.ReleaseBuffer(availableFrameSize); err != nil {
				return
			}
		}
	}

	io.Copy(audio, bytes.NewBuffer(output))

	if err = ac.Stop(); err != nil {
		return
	}
	return
}
