"""
==========
Features Visualization
==========
"""

import sys
import mpld3

import numpy as np
from numpy import genfromtxt

import matplotlib.pyplot as plt
import matplotlib.patches as mpatches
from matplotlib import colors

from mpl_toolkits.mplot3d import Axes3D
from sklearn.decomposition import PCA, KernelPCA, FastICA, IncrementalPCA
from sklearn.manifold import TSNE

print(__doc__)

x_from_dump = len(sys.argv) == 1
tsne_from_dump = True


if not tsne_from_dump:
	if x_from_dump:
		with open('./dump.pickle', 'r') as f:
			X = np.load(f)
	else:
		data_path = sys.argv[1]
		dtypes = (float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, '|S15', float)
		X = genfromtxt(data_path, delimiter=',', dtype=dtypes)
		with open('./dump.pickle', 'w+') as f:
			X.dump(f)

	print "Finished reading"

	np.random.shuffle(X)

	print "Finished shuffling"

	X = X[:20000]

	y = [seq[-1] for seq in X]
	ips = [seq[-2] for seq in X]
	X = [tuple(seq)[0:-2] for seq in X]

	print "Finished slicing and transforming"

	#kpca = KernelPCA(kernel="rbf", fit_inverse_transform=True, gamma=10)
	#X_kpca = kpca.fit_transform(X)
	#X_plot = kpca.inverse_transform(X_kpca)

	X_plot = X
	#X_plot = PCA(n_components=20).fit_transform(X)
	#X_plot = IncrementalPCA(n_components=5, batch_size=10).fit_transform(X)

	print "Finished PCA"

	#ICA
	#rng = np.random.RandomState(42)
	#ica = FastICA(random_state=rng)
	#X_plot = ica.fit(X).transform(X)  # Estimate the sources

	#t-sne
	X = X_plot
	model = TSNE(n_components=2, random_state=0)
	np.set_printoptions(suppress=True)
	X_plot = model.fit_transform(X)

	with open('./tsne_x.pickle', 'w+') as f:
		X_plot.dump(f)
	with open('./tsne_y.pickle', 'w+') as f:
		y_plot = np.asarray(y)
		y_plot.dump(f)
	with open('./ips.pickle', 'w+') as f:
		ips_list = np.asarray(ips)
		print ips_list
		ips_list.dump(f)
else:
	with open('./tsne_x.pickle', 'r') as f:
		X_plot = np.load(f)
	with open('./tsne_y.pickle', 'r') as f:
		y_plot = np.load(f)
		y = y_plot.tolist()
	with open('./ips.pickle', 'r') as f:
		ips = np.load(f)
        ips = ips.tolist()

print "Finished T-SNE"

#fig = plt.figure(2, figsize=(8, 6))

patches = []
patches.append(mpatches.Patch(color='blue', label='Unknown'))
patches.append(mpatches.Patch(color='red', label='Bot'))
patches.append(mpatches.Patch(color='green', label='Human'))
patches.append(mpatches.Patch(color='grey', label='Irrelevant'))
plt.legend(handles=patches)

# Plot the training points
cMap = colors.ListedColormap(['blue', 'red','green', 'grey'], 'indexed', 4)
bounds=[0,1,2,3,4]
norm = colors.BoundaryNorm(bounds, cMap.N)
plt.xlabel('1st eigenvector')
plt.ylabel('2nd eigenvector')

#Labels
fig, ax = plt.subplots(subplot_kw=dict(axisbg='#EEEEEE'))
fig.set_figheight(12)
fig.set_figwidth(12)

scatter = ax.scatter(X_plot[:, 0], X_plot[:, 1], c=y, cmap=cMap, norm=norm)
ax.grid(color='white', linestyle='solid')

ax.set_title("Session viz", size=20)

# Define some CSS to control our custom labels
css = """
table
{
  border-collapse: collapse;
}
th
{
  color: #ffffff;
  background-color: #000000;
}
td
{
  background-color: #cccccc;
}
table, th, td
{
  font-family:Arial, Helvetica, sans-serif;
  border: 1px solid black;
  text-align: right;
}
"""

tooltip = mpld3.plugins.PointHTMLTooltip(scatter, labels=ips, css=css)
mpld3.plugins.connect(fig, tooltip)

mpld3.show()

#
# To getter a better understanding of interaction of the dimensions
# plot the first three PCA dimensions
# fig = plt.figure(1, figsize=(8, 6))
# ax = Axes3D(fig, elev=-150, azim=110)
# ax.scatter(X_plot[:, 0], X_plot[:, 1], X_plot[:, 2], c=y, cmap=plt.cm.Paired)
# ax.set_title("Visualization")
# ax.set_xlabel("1st eigenvector")
# ax.w_xaxis.set_ticklabels([])
# ax.set_ylabel("2nd eigenvector")
# ax.w_yaxis.set_ticklabels([])
# ax.set_zlabel("3rd eigenvector")
# ax.w_zaxis.set_ticklabels([])

print "Plotted"
