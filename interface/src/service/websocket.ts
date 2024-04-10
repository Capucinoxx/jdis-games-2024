type MessageCallback = (msg: ArrayBuffer) => void;

class WebsocketService {
  private socket: WebSocket | null = null;
  private cb: MessageCallback | null = null;

  constructor() {}

  /**
   * connect permet de se connecter à un serveur websocket
   * Si une connexion est déjà établie, elle est fermée
   * et réinitialisée.
   * 
   * @param url  {string} url du serveur websocket
   */
  public connect(url: string): void {
    this.disconnect();

    this.socket = new WebSocket(url);
    this.socket.binaryType = 'arraybuffer';

    this.socket.onmessage = this.handle_message;
  }

  /**
   * disconnect permet de fermer la connexion websocket
   * si une connexion est établie.
   */
  public disconnect(): void {
    if (this.socket) {
      this.socket.close();
      this.socket = null;
    }
  }

  /**
   * subscribe permet de s'abonner à la réception de messages
   * depuis le serveur websocket.
   * @param cb {MessageCallback} callback appelée lors de la réception d'un message
   */
  public subscribe(cb : MessageCallback): void {
    this.cb = cb;
  }

  /**
   * handle_message est appelée lors de la réception d'un message
   * depuis le serveur websocket. Si un callback est défini, il est
   * appelé avec le message reçu lorsque celui-ci est un ArrayBuffer.
   * 
   * @param event {MessageEvent} événement de réception de message
   */
  private handle_message = (event: MessageEvent): void => {
    if (this.cb && event.data instanceof ArrayBuffer)
      this.cb(event.data);
  }
};

export { WebsocketService };