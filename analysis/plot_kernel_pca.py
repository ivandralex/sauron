"""
==========
Features Visualization
==========
"""

import sys

import numpy as np
from numpy import genfromtxt

import matplotlib.pyplot as plt
import matplotlib.patches as mpatches
from matplotlib.colors import ListedColormap

from mpl_toolkits.mplot3d import Axes3D
from sklearn.decomposition import PCA, KernelPCA, FastICA, IncrementalPCA
from sklearn.manifold import TSNE

print(__doc__)

from_dump = len(sys.argv) == 1

if from_dump:
    with open('./dump.pickle', 'r') as f:
        X = np.load(f)
else:
    data_path = sys.argv[1]
    X = genfromtxt(data_path, delimiter=',')
    with open('./dump.pickle', 'w+') as f:
        X.dump(f)

print "Finished reading"

np.random.shuffle(X)

print "Finished shuffling"

#X = X[:10000,:]
X = X[:1000,:]

#Counts only
num_features = X.shape[1]
y = X[:, num_features - 1 : num_features]
X = X[:, :num_features]

print X.shape
print y.shape

y = y.tolist()

print "Finished slicing and transforming"


#kpca = KernelPCA(kernel="rbf", fit_inverse_transform=True, gamma=10)
#X_kpca = kpca.fit_transform(X)
#X_plot = kpca.inverse_transform(X_kpca)

#gX_plot = PCA(n_components=3).fit_transform(X)
X_plot = IncrementalPCA(n_components=5, batch_size=10).fit_transform(X)

#ICA
#rng = np.random.RandomState(42)
#ica = FastICA(random_state=rng)
#X_plot = ica.fit(X).transform(X)  # Estimate the sources

#t-sne
X = X_plot
model = TSNE(n_components=2, random_state=0)
np.set_printoptions(suppress=True)
X_plot = model.fit_transform(X)

plt.figure(2, figsize=(8, 6))

patches = []
patches.append(mpatches.Patch(color='blue', label='Unknown'))
patches.append(mpatches.Patch(color='red', label='Bot'))
patches.append(mpatches.Patch(color='green', label='Human'))
patches.append(mpatches.Patch(color='grey', label='Irrelevant'))
plt.legend(handles=patches)

# Plot the training points
cMap = ListedColormap(['blue', 'red','green', 'grey'])
plt.scatter(X_plot[:, 0], X_plot[:, 1], c=y, cmap=cMap)
plt.xlabel('1st eigenvector')
plt.ylabel('2nd eigenvector')
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

plt.show()
