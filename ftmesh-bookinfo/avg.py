def read_numbers(filename):
    # 初始化列表，用于存储读取的数字
    numbers = []

    # 打开文件并读取每行数据
    with open(filename, 'r') as file:
        for line in file:
            # 将每行的内容转换为浮点数并添加到列表中
            numbers.append(float(line.strip()))

    return numbers

def calculate_average(numbers):
    # 计算平均值
    if numbers:  # 检查列表不为空
        return sum(numbers) / len(numbers)
    else:
        return 0

# 主程序
filename = 'num.txt'
numbers = read_numbers(filename)  # 读取数据
average = calculate_average(numbers)  # 计算平均值

print(f"The average of the numbers is: {average:.2f}")
