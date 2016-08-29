"""
==========
Features Visualization
==========
"""

import numpy as np
import matplotlib.pyplot as plt
from mpl_toolkits.mplot3d import Axes3D
from sklearn.decomposition import PCA, KernelPCA, FastICA, IncrementalPCA
from sklearn.manifold import TSNE

print(__doc__)

from numpy import genfromtxt
X = genfromtxt('../output/features/21_10k.csv', delimiter=',')

#Original data with label
#y = X[:, 105:106]
#X = X[:, :105]

#Select path vector
#5n - 4
#X = X[:, 51:55]

#Counts onlys
y = X[:, 21:22]
X = X[:, :21]

print X.shape
print y.shape

y = y.tolist()

#kpca = KernelPCA(kernel="rbf", fit_inverse_transform=True, gamma=10)
#X_kpca = kpca.fit_transform(X)
#X_plot = kpca.inverse_transform(X_kpca)

#gX_plot = PCA(n_components=3).fit_transform(X)
X_plot = IncrementalPCA(n_components=3, batch_size=10).fit_transform(X)

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

# Plot the training points
plt.scatter(X_plot[:, 0], X_plot[:, 1], c=y, cmap=plt.cm.Paired)
plt.xlabel('1st eigenvector')
plt.ylabel('2nd eigenvector')

# To getter a better understanding of interaction of the dimensions
# plot the first three PCA dimensions
#fig = plt.figure(1, figsize=(8, 6))
#ax = Axes3D(fig, elev=-150, azim=110)
#ax.scatter(X_plot[:, 0], X_plot[:, 1], X_plot[:, 2], c=y, cmap=plt.cm.Paired)
#ax.set_title("First three PCA directions")
#ax.set_xlabel("1st eigenvector")
#ax.w_xaxis.set_ticklabels([])
#ax.set_ylabel("2nd eigenvector")
#ax.w_yaxis.set_ticklabels([])
#ax.set_zlabel("3rd eigenvector")
#ax.w_zaxis.set_ticklabels([])

plt.show()
