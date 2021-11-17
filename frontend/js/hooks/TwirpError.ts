export class TwirpError extends Error {
  static fromJson(json: string): TwirpError {
    // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
    const o = JSON.parse(json);
    return new TwirpError(o.code, o.msg);
  }

  get message(): string {
    return `${this.code}:${this.msg}`;
  }

  constructor(code: string, message: string) {
    super();
    this.code = code;
    this.msg = message;
  }

  code: string;
  msg: string;
}
