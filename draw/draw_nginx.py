import matplotlib.pyplot as plt
import numpy as np

# %% Style
SECONDARY_COLOR = "#808080"
FONT_SIZE_SM = 18
FONT_SIZE_MD = 20
FONT_SIZE_LG = 22
COLORS = [
    "tab:blue",
    "tab:orange",
    "tab:green",
    "tab:red",
    "tab:purple",
    "tab:brown",
    "tab:pink",
    "tab:gray",
    "tab:olive",
    "tab:cyan",
]
HATCHS = ["--", "//", "xx", "\\", "oo", "..", "+", "O", "*"]

plt.rcParams.update(
    {
        "figure.figsize": (10, 6),
        "font.size": FONT_SIZE_SM,
        "axes.titlesize": FONT_SIZE_LG,
        "axes.labelsize": FONT_SIZE_LG,
        "font.family": "Calibri",
        "legend.framealpha": 0.5,
        "legend.fontsize": FONT_SIZE_SM,
        "axes.spines.right": False,
        "axes.spines.top": False,
        "axes.edgecolor": SECONDARY_COLOR,
        "xtick.color": SECONDARY_COLOR,
        "ytick.color": SECONDARY_COLOR,
        "xtick.labelsize": FONT_SIZE_LG,
        "ytick.labelsize": FONT_SIZE_LG,
        "grid.color": SECONDARY_COLOR,
        "grid.linestyle": ":",
        "grid.linewidth": 1,
        # enable grid
        "axes.grid": True,
        # disable x axis grid
        "axes.grid.axis": "y",
    }
)

if __name__ == "__main__":
    fig = plt.figure()
    ax = fig.add_subplot()

    x_labels = ["0.5KB+0.5KB", "1.5KB+0.5KB", "10KB+0.5KB","10KB+10KB"]
    # company_y = [0.0,0.42857]
    # bics_y = [0.141558,0.17316]
    # ft4_y = [0, 0.1]
    ydata = [
        [1539.56, 1550.48, 1589.36,1589.36, "NoSync"],
        [1674.20, 1682.92, 1774.08,1924.80, "MeshSync"],
        [4462.08, 4540.48, 4600.32,4708.98, "SepSync"],
        # [0.04523809524,0.05428571429, 1, 'NetE'],
        # [0.0582278481,0.05907172996, 1, 'NetF'],
        # [0,0, 1, 'NetG'],
        # [0,0, 1, 'NetH'],
        # [1,1, 1, 'NetI'],
    ]
    
    for row in ydata:
        row[:len(x_labels)] = [value / 1000 for value in row[:len(x_labels)]]

    #ax.grid(True, alpha=0.5)
    ax.grid(True, alpha=0.5, which='major', axis='y')
    for idx in range(0, len(ydata)):
        hdl = ax.bar(
            [i - 0.36 + idx * (0.95 / len(ydata)) for i in range(0, len(x_labels))],
            ydata[idx][0 : len(x_labels)],
            align="center",
            hatch=HATCHS[idx],
            color="none",
            edgecolor=COLORS[idx],
            width=0.14,
            label=ydata[idx][-1],
        )

    ax.set_xticks(np.arange(len(x_labels)))
    ax.set_xticklabels(x_labels, fontsize=20)
    ax.set_ylabel("Response Latency (ms)")
    # ax.set_ylim(0, 1.1)
    #ax.set_yscale("log")

    ax.legend(framealpha=0.3, prop={"size": 20}, draggable=True, ncols=3, labelspacing=0.2,loc='upper center',bbox_to_anchor=(0.5, 1.15))
    plt.tight_layout()
    plt.savefig('nginx_overhead.png', format='png', dpi=400)
    plt.show()