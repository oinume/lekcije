from __future__ import print_function
from apscheduler.schedulers.blocking import BlockingScheduler
import logging
import subprocess

logging.basicConfig()
job_defaults = {
    'coalesce': False,
    'max_instances': 2
}
scheduler = BlockingScheduler(job_defaults=job_defaults)


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


@scheduler.scheduled_job('interval', days=7)
def timed_job_days7():
    print("Run teacher_error_resetter")
    subprocess.check_call(
        "teacher_error_resetter -concurrency=5",
        shell=True)

scheduler.start()
