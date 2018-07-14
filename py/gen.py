import sys
import random

n = int(sys.argv[1])

for i in range(n):
    key = ' '.join([ str(random.randint(0, 15)) for _ in range(8) ])
    value = str(random.randint(0,10e6))
    print(key)
    print(value)
