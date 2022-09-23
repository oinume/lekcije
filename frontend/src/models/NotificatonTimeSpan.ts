// TODO: Add method isZero() and parse().
type MaybeNotificationTimeSpan = {
  fromHour: number;
  fromMinute: number;
  toHour: number;
  toMinute: number;
};

export class NotificationTimeSpanModel {
  [key: string]: any

  static fromObject(o: MaybeNotificationTimeSpan): NotificationTimeSpanModel {
    return new NotificationTimeSpanModel(o.fromHour, o.fromMinute, o.toHour, o.toMinute);
  }

  constructor(
    public fromHour: number,
    public fromMinute: number,
    public toHour: number,
    public toMinute: number,
  ) {}

  isZero(): boolean {
    return this.fromHour === 0 && this.fromMinute === 0 && this.toHour === 0 && this.toMinute === 0;
  }
}
