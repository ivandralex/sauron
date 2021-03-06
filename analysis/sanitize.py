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

data_path = sys.argv[1]
output_path = sys.argv[2]
label = sys.argv[3]
data = pandas.read_csv(data_path, delimiter=';')

data = data.dropna()
ips = data['ip'].values

print len(ips)

with open(output_path, 'w+') as f:
    for ip in ips:
        print ip
        f.write(ip + '\n')
    f.close()
