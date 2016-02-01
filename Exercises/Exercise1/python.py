# Bruk Ideone.com for å kompilere eller http://www.tutorialspoint.com/execute_python_online.php


from threading import Thread

i = 0

def someThreadFunction1():    
# Potentially useful thing:
#   In Python you "import" a global variable, instead of "export"ing it when you declare it
#   (This is probably an effort to make you feel bad about typing the word "global")
    global i
    for j in range (0,1000000):
        i += 1
def someThreadFunction2():
    global i
    for j in range (0,1000000):
        i -= 1
    
def main():
    someThread1 = Thread(target = someThreadFunction1, args = (),)
    someThread1.start()
    someThread2 = Thread(target = someThreadFunction2, args = (),)
    someThread2.start()
    
    someThread1.join()
    someThread2.join()
    print("Did it work?")
    print(i)

main()

#Gir forskjellige svar ulike null fordi threadene kjøres samtidig (tokjernet) (?)
