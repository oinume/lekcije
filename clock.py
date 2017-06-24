from apscheduler.schedulers.blocking import BlockingScheduler
import subprocess

scheduler = BlockingScheduler()


@scheduler.scheduled_job('interval', minutes=1)
def timed_job():
    print("Run notifier")
    subprocess.run(
        "notifier -concurrency=5 -fetcher-cache=true -notification-interval=1 && curl -sS https://nosnch.in/c411a3a685",
        shell=True,
        check=True)

scheduler.start()
