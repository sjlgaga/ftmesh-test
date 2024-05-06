import matplotlib.pyplot as plt

# 读取数据
def read_data(filename):
    with open(filename, 'r') as file:
        lines = file.readlines()
    
    headers = lines[0].strip().split()  # 第一行为标题，包含所有元素
    labels = []
    values = []
    for line in lines[1:]:  # 从第二行开始是数据
        parts = line.split()
        labels.append(parts[0])  # 提取每一行的第一列作为x轴标签
        values.append([round(float(part) / 1000, 2) for part in parts[1:]])  # 将剩下的数值转换为浮点数并保存
    
    return labels, values, headers[1:]  # 返回标签，数值和除了第一个元素外的头部

# 绘制直方图
def plot_histogram(labels, values, headers):
    # 转置values以便每一列成为一个柱状组
    transposed_values = list(zip(*values))
    
    x = range(len(labels))  # x轴的位置
    width = 0.25  # 柱子的宽度，调整以适应三组数据

    fig, ax = plt.subplots()  # 创建一个图和坐标轴

    # 为每个数据集绘制柱状图
    for i, header in enumerate(headers):
        bars = ax.bar([p + width*i for p in x], transposed_values[i], width=width, label=header, align='center')
    
    # 设置图例和标签
    ax.set_xlabel('Request Size(+ attached states size) (Byte)')
    ax.set_ylabel('Latency (ms)')
    ax.set_title('Performance Metrics')
    ax.set_xticks([p + width for p in x])  # 设置x轴标签为组的中心
    ax.set_xticklabels(labels)
    ax.legend()  # 显示图例
    
    # 调整Y轴范围，增加上限20%
    ax.set_ylim(0, max(max(value) for value in values) * 1.4)

    # 保存图表为PNG文件
    plt.savefig('performance_metrics.png', dpi=400)
    
    # 显示图表
    plt.show()

# 主程序
filename = 'data.txt'
labels, values, headers = read_data(filename)
plot_histogram(labels, values, headers)
