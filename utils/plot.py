import matplotlib.pyplot as plt


f = open("output/output.txt", "r")

data = []
for x in f.readlines():
    splitData = x.split()
    data.append(splitData)

for i in range(0,10):
    ypoints = [int(x[1]) for x in data if x[0] == str(i)]
    xpoints = [x for x in range(0,len(ypoints))]
    plt.plot(xpoints, ypoints)

plt.show()
