"""
==========
Features Visualization
==========
"""

import sys
import string
import mpld3

import numpy as np
import pandas
from numpy import genfromtxt

import matplotlib.pyplot as plt
import matplotlib.patches as mpatches
from matplotlib import colors

from scipy.sparse import csr_matrix

from mpl_toolkits.mplot3d import Axes3D
from sklearn.decomposition import PCA, KernelPCA, FastICA, IncrementalPCA, TruncatedSVD
from sklearn.manifold import TSNE

print(__doc__)

x_from_dump = len(sys.argv) == 1
tsne_from_dump = False #en(sys.argv) == 1


if not tsne_from_dump:
	if x_from_dump:
		data = pandas.read_pickle('./dumps/dump.pickle')
	else:
		data_path = sys.argv[1]
		data = pandas.read_csv(data_path)
		data.to_pickle('./dumps/dump.pickle')

	print "Finished reading"

	data = data[:200000]
	#np.random.shuffle(X)
	#print "Finished shuffling"

	y = data.iloc[:, -1:]
	X = data.iloc[:, 2:-1]

	print X.shape
	print y.shape

	print "Rows: %s" % len(X)

	#Keys
	ips = []
	for seq in data.values:
		ips.append(str(seq[0]) + "|" + str(seq[1]))

	data = data.drop('user_agent', 1)
	data = data.drop('ip', 1)

	print "Finished slicing and transforming"

	n_components = 2

	model = 'pca'

	if model == 'pca':
		pca = PCA(n_components=n_components)
		X_plot = pca.fit_transform(X)
	if model == 'kernel-pca':
		kpca = KernelPCA(kernel="rbf", fit_inverse_transform=True, gamma=10)
		X_kpca = kpca.fit_transform(X)
		X_plot = kpca.inverse_transform(X_kpca)
	if model == 'truncated-svd':
		X_sparse = csr_matrix(X)
		pca = TruncatedSVD(n_components=n_components, random_state=241)
		X_plot = pca.fit_transform(X_sparse)
	if model == 'incremental-pca':
		pca = IncrementalPCA(n_components=n_components, batch_size=10, copy=False)
		X_plot = pca.fit_transform(X)
	if model == 'fast-ica':
		rng = np.random.RandomState(42)
		ica = FastICA(random_state=rng)
		X_plot = ica.fit(X).transform(X)  # Estimate the sources

	if 'pca' in model or 'svd' in model:
		print pca.explained_variance_ratio_

		dominator_features = 10

		print "Finished PCA: %.3f of variance retained" % np.sum(pca.explained_variance_ratio_)
		i = 0
		while i<n_components:
			print "Feature %s dominated by: %s\n" % (i, data.columns.values[np.argpartition(pca.components_[0], -dominator_features)[-dominator_features:]])
			i = i + 1
		print("\n\n\n\n")

	#t-sne
	if model == 't-sne':
		model = TSNE(n_components=n_components, random_state=241)
		np.set_printoptions(suppress=True)
		X_plot = model.fit_transform(X)

	#Dump features
	with open('./dumps/tsne_x.pickle', 'w+') as f:
		X_plot.dump(f)
	with open('./dumps/tsne_y.pickle', 'w+') as f:
		y_plot = np.asarray(y)
		y_plot.dump(f)
	with open('./dumps/ips.pickle', 'w+') as f:
		ips_list = np.asarray(ips)
		print ips_list
		ips_list.dump(f)
else:
	with open('./dumps/tsne_x.pickle', 'r') as f:
		X_plot = np.load(f)
	with open('./dumps/tsne_y.pickle', 'r') as f:
		y_plot = np.load(f)
		y = y_plot.tolist()
	with open('./dumps/ips.pickle', 'r') as f:
		ips = np.load(f)
        ips = ips.tolist()

print "Finished T-SNE"

# Plot the training points
cMap = colors.ListedColormap(['blue', 'red','green', 'grey'], 'indexed', 4)
bounds=[0,1,2,3,4]
norm = colors.BoundaryNorm(bounds, cMap.N)
plt.xlabel('1st eigenvector')
plt.ylabel('2nd eigenvector')

#Labels
fig, ax = plt.subplots(subplot_kw=dict(axisbg='#EEEEEE'))
fig.set_figheight(14)
fig.set_figwidth(24)

scatter = ax.scatter(X_plot[:, 0], X_plot[:, 1], c=y, cmap=cMap, norm=norm)
ax.grid(color='white', linestyle='solid')

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

#CLICK
class ClickInfo(mpld3.plugins.PluginBase):
    """Plugin for getting info on click"""

    #TODO: move this to separate js
    JAVASCRIPT = "window.ips = [\"" + string.join(ips, "\",\"") + "\"]"
    JAVASCRIPT += """
    mpld3.register_plugin("clickinfo", ClickInfo);
    ClickInfo.prototype = Object.create(mpld3.Plugin.prototype);
    ClickInfo.prototype.constructor = ClickInfo;
    ClickInfo.prototype.requiredProps = ["id"];
    function ClickInfo(fig, props){
        mpld3.Plugin.call(this, fig, props);
    };

    ClickInfo.prototype.draw = function(){
            var obj = mpld3.get_element(this.props.id);

            obj.elements().on("mousedown", function(d, i){
                var ip = ips[Number(i)];
                console.log(ip);
                var el = document.getElementById('ipLabel');
                if(!el){
                    var el = document.createElement('p')
                    el.id = 'ipLabel';
                    el.style.cssText = 'position: absolute; top: 0; left: 20';
                    document.body.appendChild(el)
                }
				var pieces = ip.split('|')
				var command = 'grep ' + pieces[0] + ' ~/repos/data-miner-utils/data/dump.list | grep "' + pieces[1] + '"';
                //el.innerHTML = ip;
				el.innerHTML = command;
                range = document.createRange();
                range.selectNode(el);
                window.getSelection().addRange(range);
                document.execCommand('copy')
        });
    }
    """
    def __init__(self, points):
        self.dict_ = {"type": "clickinfo",
                      "id": mpld3.utils.get_id(points),}

mpld3.plugins.connect(fig, ClickInfo(scatter))

#~CLICK

tooltip = mpld3.plugins.PointHTMLTooltip(scatter, labels=ips, css=css)
mpld3.plugins.connect(fig, tooltip)

mpld3.show()

print "Plotted"
