export class FormErrors<T> {
  message: string;
  errors: Record<keyof T, string>;

  constructor(message: string, error: Record<keyof T, string>) {
    this.errors = error;
    this.message = message;
  }
}
