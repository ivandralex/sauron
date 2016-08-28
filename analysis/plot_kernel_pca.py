"""
==========
Features Visualization
==========
"""

import numpy as np
import matplotlib.pyplot as plt
from mpl_toolkits.mplot3d import Axes3D
from sklearn.decomposition import PCA, KernelPCA, FastICA, IncrementalPCA

print(__doc__)

from numpy import genfromtxt
X = genfromtxt('../output/features/100000_with_labels.csv', delimiter=',')

#Original data with label
y = X[:, 105:106]
X = X[:, :105]

#Select path vector
#5n - 4
X = X[:, 51:55]

print X.shape
print y.shape

print y[:1000]

#kpca = KernelPCA(kernel="rbf", fit_inverse_transform=True, gamma=10)
#X_kpca = kpca.fit_transform(X)
#X_plot = kpca.inverse_transform(X_kpca)

#X_plot = PCA(n_components=3).fit_transform(X)
X_plot = IncrementalPCA(n_components=3, batch_size=10).fit_transform(X)

#ICA
#rng = np.random.RandomState(42)
#ica = FastICA(random_state=rng)
#X_plot = ica.fit(X).transform(X)  # Estimate the sources

plt.figure(2, figsize=(8, 6))
plt.clf()

# Plot the training points
plt.scatter(X_plot[:, 1], X_plot[:, 2], c=y, cmap=plt.cm.Paired)
plt.xlabel('1st eigenvector')
plt.ylabel('2nd eigenvector')

# To getter a better understanding of interaction of the dimensions
# plot the first three PCA dimensions
fig = plt.figure(1, figsize=(8, 6))
ax = Axes3D(fig, elev=-150, azim=110)
ax.scatter(X_plot[:, 0], X_plot[:, 1], X_plot   [:, 2], c=y, cmap=plt.cm.Paired)
ax.set_title("First three PCA directions")
ax.set_xlabel("1st eigenvector")
ax.w_xaxis.set_ticklabels([])
ax.set_ylabel("2nd eigenvector")
ax.w_yaxis.set_ticklabels([])
ax.set_zlabel("3rd eigenvector")
ax.w_zaxis.set_ticklabels([])

plt.show()
