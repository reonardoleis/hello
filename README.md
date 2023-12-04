![hello](https://i.imgur.com/oTgyymW.png)

# hello

Hello is a simple chat service written in Go that allows users to communicate over TCP. It provides basic chat functionality and aims to evolve into a more feature-rich chat platform in the future.

## Features

- **Basic Chat:** Users can connect to the Hello service using a TCP client and engage in real-time chat with others.

## Future Roadmap

The Hello service is actively under development, and we have exciting features planned for the future releases:

- **Chat Rooms:** Introduce the concept of chat rooms for better organization and focused discussions.

- **Password Protection:** Enhance the security of chat rooms by implementing password protection.

- **Member List:** Display a list of members within a chat room for better visibility and user management.

- **Voice Chat:** Explore the integration of voice chat capabilities to enable richer communication experiences.

## Getting Started

To run the Hello service locally, follow these steps:

1. Clone the repository: `git clone https://github.com/your-username/hello.git`
2. Navigate to the project directory: `cd hello`
3. Build the executable: `make build-server && make build-client`
4. Run the server: `./bin/hello-server <PORT>`
5. Run the client: `./bin/hello-client <HOST:PORT>`

Now, the Hello service should be running locally, and clients can connect to it using a TCP client.
