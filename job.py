#!/usr/bin/python3
import os, random, subprocess, time


def getAndSet():
    cwd = os.getcwd() + "/image"
    imgName = random.choice(os.listdir(cwd))
    fullPath = cwd + '/' + imgName
    command = 'gsettings set org.gnome.desktop.background picture-uri ' + 'file://' + fullPath
    print("full command is:", command)
    process = subprocess.Popen(command.split(), stdout=subprocess.PIPE)
    output, error = process.communicate()


def job(x: int):
    while True:
        getAndSet()
        time.sleep(60 * x)


if __name__ == "__main__":
    job(1)
