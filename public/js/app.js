class WSClient {
    /** @type {WebSocket} */
    #ws = null;

    constructor() {
        this.#ws = new WebSocket('ws://localhost:8080/ws');
        this.#ws.onopen = () => {
            console.log('Connected to server');
        };

        this.#ws.onclose = () => {
            console.log('Disconnected from server');
        }

        this.#ws.onmessage = (message) => {
            console.log('Received message:', message.data);
        };
    }

    createRoom(roomName) {
        this.#ws.send(JSON.stringify({
            type: 'command-create-room',
            data: {
                roomName
            }
        }));
    }

    joinRoom(roomId) {
        this.#ws.send(JSON.stringify({
            type: 'command-join-room',
            data: {
                roomId
            }
        }));
    }

    sendMessage(roomId, message) {
        this.#ws.send(JSON.stringify({
            type: 'command-send-message',
            data: {
                roomId,
                message
            }
        }));
    }

    leaveRoom(roomId) {
        this.#ws.send(JSON.stringify({
            type: 'command-leave-room',
            data: {
                roomId
            }
        }));
    }
}

