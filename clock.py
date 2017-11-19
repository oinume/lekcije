from __future__ import print_function
from apscheduler.schedulers.blocking import BlockingScheduler
import logging
import subprocess

logging.basicConfig()
scheduler = BlockingScheduler()


@scheduler.scheduled_job('interval', minutes=1)
def timed_job_min1():
    print("Run notifier (interval=1)")
    subprocess.check_call(
        "notifier -concurrency=5 -fetcher-cache=true -notification-interval=1 && curl -sS https://nosnch.in/c411a3a685",
        shell=True)


@scheduler.scheduled_job('interval', minutes=10)
def timed_job_min10():
    print("Run notifier (interval=10)")
    subprocess.check_call(
        "notifier -concurrency=9 -fetcher-cache=true -notification-interval=10 && curl -sS https://nosnch.in/c411a3a685",
        shell=True)

scheduler.start()
