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

from mpl_toolkits.mplot3d import Axes3D
from sklearn.decomposition import PCA, KernelPCA, FastICA, IncrementalPCA
from sklearn.manifold import TSNE

print(__doc__)

x_from_dump = len(sys.argv) == 1
tsne_from_dump = len(sys.argv) == 1


if not tsne_from_dump:
	if x_from_dump:
		with open('./dump.pickle', 'r') as f:
			X = np.load(f)
	else:
		data_path = sys.argv[1]
		data = pandas.read_csv(data_path)
		X = data.as_matrix()# genfromtxt(data_path, delimiter=',',usecols=np.arange(0,1434))
		with open('./dump.pickle', 'w+') as f:
			X.dump(f)

	print "Finished reading"

	np.random.shuffle(X)

	print "Rows: %s" % len(X)

	#X = X[:110000]

	print "Finished reading"

	np.random.shuffle(X)

	print "Finished shuffling"

	y = [seq[-1] for seq in X]
	keys = [seq[0:2] ]
	ips = []
	for seq in X:
		ips.append(seq[0] + "|" + seq[1])
	X = [tuple(seq)[2:-1] for seq in X]

	data = data.drop('user_agent', 1)
	data = data.drop('ip', 1)

	print "Finished slicing and transforming"

	#kpca = KernelPCA(kernel="rbf", fit_inverse_transform=True, gamma=10)
	#X_kpca = kpca.fit_transform(X)
	#X_plot = kpca.inverse_transform(X_kpca)

	n_components = 2

	pca = PCA(n_components=n_components)
	X_plot = pca.fit_transform(X[:])

	print pca.explained_variance_ratio_

	#X_plot = IncrementalPCA(n_components=2, batch_size=10).fit_transform(X)
	dominator_features = 10

	print "Finished PCA: %.3f of variance retained" % np.sum(pca.explained_variance_ratio_)
	i = 0
	while i<n_components:
		print "Feature %s dominated by: %s\n" % (i, data.columns.values[np.argpartition(pca.components_[0], -dominator_features)[-dominator_features:]])
		i = i + 1

	print("\n\n\n\n")

	sys.exit(0)

	#ICA
	#rng = np.random.RandomState(42)
	#ica = FastICA(random_state=rng)
	#X_plot = ica.fit(X).transform(X)  # Estimate the sources

	#t-sne
	#model = TSNE(n_components=2, random_state=241)
	#np.set_printoptions(suppress=True)
	#X_plot = model.fit_transform(X_plot)

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
                el.innerHTML = ip;
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
