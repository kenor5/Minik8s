import time

print("Starting CPU load...")

while True:
    # 计算一些无用操作
    for i in range(10000000):
        x = i * i
        y = x + i

    # 停顿一段时间
    time.sleep(0.5)
