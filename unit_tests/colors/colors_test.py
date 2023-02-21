import sys
for i in range(0, 16):
    for j in range(0, 16):
        code = str(i * 16 + j)
        sys.stdout.write("\033[48;5;" + code + "m " + code.ljust(4))

        sys.stdout.write(code)
    print(u"\u001b[0m")


print("\033[3;2m")
print("TESTING")
print(u"\u001b[0m")