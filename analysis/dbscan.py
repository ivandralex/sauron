"""
==========
Features Visualization
==========
"""

import sys
import string
import mpld3

import numpy as np
from numpy import genfromtxt

import matplotlib.pyplot as plt
import matplotlib.patches as mpatches
from matplotlib import colors

from sklearn.cluster import DBSCAN
from sklearn import metrics
from sklearn.datasets.samples_generator import make_blobs
from sklearn.preprocessing import StandardScaler

from mpl_toolkits.mplot3d import Axes3D
from sklearn.decomposition import PCA, KernelPCA, FastICA, IncrementalPCA
from sklearn.manifold import TSNE

print(__doc__)

x_from_dump = len(sys.argv) == 1
tsne_from_dump = False#len(sys.argv) == 1


if not tsne_from_dump:
	if x_from_dump:
		with open('./dump.pickle', 'r') as f:
			X = np.load(f)
	else:
		data_path = sys.argv[1]
		dtypes = ('|S15', float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float, float)
		X = genfromtxt(data_path, delimiter=',', dtype=dtypes)
		with open('./dump.pickle', 'w+') as f:
			X.dump(f)

	print "Finished reading"

	np.random.shuffle(X)

	print "Finished shuffling"

	X = X[:20000]

	y = [seq[-1] for seq in X]
	ips = [seq[0] for seq in X]
	X = [tuple(seq)[1:-1] for seq in X]

	print "Finished slicing and transforming"

	X_plot = X

else:
	with open('./tsne_x.pickle', 'r') as f:
		X_plot = np.load(f)
	with open('./tsne_y.pickle', 'r') as f:
		y_plot = np.load(f)
		y = y_plot.tolist()
	with open('./ips.pickle', 'r') as f:
		ips = np.load(f)
        ips = ips.tolist()


labels_true = y
X = np.array(X_plot)

#from sklearn.ensemble import RandomForestClassifier
#clf = RandomForestClassifier(n_estimators=10)
#clf = clf.fit(X, Y)

db = DBSCAN(eps=0.3, min_samples=10, n_jobs=4).fit(X)
core_samples_mask = np.zeros_like(db.labels_, dtype=bool)
core_samples_mask[db.core_sample_indices_] = True
labels = db.labels_

# Number of clusters in labels, ignoring noise if present.
n_clusters_ = len(set(labels)) - (1 if -1 in labels else 0)

print('Estimated number of clusters: %d' % n_clusters_)
print("Homogeneity: %0.3f" % metrics.homogeneity_score(labels_true, labels))
print("Completeness: %0.3f" % metrics.completeness_score(labels_true, labels))
print("V-measure: %0.3f" % metrics.v_measure_score(labels_true, labels))
print("Adjusted Rand Index: %0.3f"
      % metrics.adjusted_rand_score(labels_true, labels))
print("Adjusted Mutual Information: %0.3f"
      % metrics.adjusted_mutual_info_score(labels_true, labels))
print("Silhouette Coefficient: %0.3f"
      % metrics.silhouette_score(X, labels))

fig = plt.figure(2, figsize=(8, 6))

patches = []
patches.append(mpatches.Patch(color='blue', label='Unknown'))
patches.append(mpatches.Patch(color='red', label='Bot'))
patches.append(mpatches.Patch(color='green', label='Human'))
patches.append(mpatches.Patch(color='grey', label='Irrelevant'))
plt.legend(handles=patches)

cMap = colors.ListedColormap(['blue', 'red','green', 'grey'], 'indexed', 4)
bounds=[0,1,2,3,4]
norm = colors.BoundaryNorm(bounds, cMap.N)

# Black removed and is used for noise instead.
unique_labels = set(labels)
colors = plt.cm.Spectral(np.linspace(0, 1, len(unique_labels)))
for k, col in zip(unique_labels, colors):
    if k == -1:
        # Black used for noise.
        col = 'k'

    class_member_mask = (labels == k)

    xy = X[class_member_mask & core_samples_mask]
    plt.plot(xy[:, 0], xy[:, 1], 'o', markerfacecolor=col,
             markeredgecolor='k', markersize=14)

    xy = X[class_member_mask & ~core_samples_mask]
    plt.plot(xy[:, 0], xy[:, 1], 'o', markerfacecolor=col,
             markeredgecolor='k', markersize=6)

plt.title('Estimated number of clusters: %d' % n_clusters_)
plt.show()


plt.show()


print "Plotted"
