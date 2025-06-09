declare module "negotiator" {
  interface Headers {
    [key: string]: string;
  }

  class Negotiator {
    constructor(headers: Record<string, string>);
    languages(): string[];
  }

  export = Negotiator;
}
