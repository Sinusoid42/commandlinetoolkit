import sys
for i in range(0, 16):
    for j in range(0, 16):
        code = str(i * 16 + j)
        sys.stdout.write("\033[48;5;" + code + "m " + code.ljust(4))

        sys.stdout.write(code)
    print(u"\u001b[0m")

for i in range(0, 16):
    for j in range(0, 16):
        code = str(i * 16 + j)
        sys.stdout.write("\033[48;5;" + code + "m " + code.ljust(4))

        sys.stdout.write(code)
    print(u"\u001b[0m")


print("\033[3;2m")
print("TESTING")
print(u"\u001b[0m")


print("\033[36;16;22m", "TEST")
print("\033[37;16;22m", "TEST")
print("\033[38;16;22m", "TEST")
print("\033[39;16;22m", "TEST")
print("\033[46;16;22m", "TEST")
print("\033[47;16;22m", "TEST")
print("\033[48;16;22m", "TEST")
print("\033[49;16;22m", "TEST")
print("\033[36;16;22m", "TEST")
print("\033[36;26;22m", "TEST")
print("\033[36;36;22m", "TEST")
print("\033[36;46;22m", "TEST")
print("\033[36;56;22m", "TEST")
print("\033[36;66;22m", "TEST")
print("\033[49;76;22m", "TEST")
print("\033[49;86;22m", "TEST")
print("\033[49;96;22m", "TEST")
print("\033[49;16;22m", "TEST")
