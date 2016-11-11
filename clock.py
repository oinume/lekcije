from apscheduler.schedulers.blocking import BlockingScheduler
import subprocess

scheduler = BlockingScheduler()

@scheduler.scheduled_job('interval', minutes=5)
def timed_job():
    print("Run notifier")
    subprocess.run("notifier && curl https://nosnch.in/c411a3a685", shell=True, check=True)

scheduler.start()
