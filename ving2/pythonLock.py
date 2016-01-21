

from threading import Thread
from threading import Lock
i = 0

def someThreadFunction1(lock):    
# Potentially useful thing:
#   In Python you "import" a global variable, instead of "export"ing it when you declare it
#   (This is probably an effort to make you feel bad about typing the word "global")
    global i
    for j in range (0,1000000):
        lock.acquire()
        i += 1
        lock.release()


def someThreadFunction2(lock):
    global i
    for j in range (0,1000000):
        lock.acquire()
        i -= 1
        lock.release()
    
def main():
    lock = Lock()
    someThread1 = Thread(target = someThreadFunction1, args = ([lock]))
    someThread1.start()
    someThread2 = Thread(target = someThreadFunction2, args = ([lock]))
    someThread2.start()
    
    someThread1.join()
    someThread2.join()
    print(i)

main()
