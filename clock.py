from apscheduler.schedulers.blocking import BlockingScheduler

scheduler = BlockingScheduler()

@scheduler.scheduled_job('interval', minutes=30)
def timed_job():
    print('This job is run every 30 minutes.')

scheduler.start()
