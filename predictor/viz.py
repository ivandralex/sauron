from keras.layers import Input, LSTM, RepeatVector
from keras.models import Model
from keras.preprocessing import sequence

import sys
import pandas
import numpy as np

timesteps = 1000
input_dim = 1000
latent_dim = 32

nb_epoch = 50
batch_size = 100

#Reading input
if len(sys.argv) > 1 :
    data_path = sys.argv[1]
    data = pandas.read_csv(data_path)
    data = data.drop('user_agent', 1)
    data = data.drop('ip', 1)
    data = data.drop('label', 1)

    X = data.as_matrix();
    np.random.shuffle(X)
    X_train = X[0:1000]
    X_test = X[1000:2000]

    with open('./train.pickle', 'w+') as f:
        X_train.dump(f)
    with open('./test.pickle', 'w+') as f:
        X_test.dump(f)
else:
    print('From dump')
    with open('./train.pickle', 'r') as f:
        X_train = np.load(f)
    with open('./test.pickle', 'r') as f:
        X_test = np.load(f)

print X_train.shape
print X_test.shape

X_train = np.reshape(X_train, X_train.shape + (1,))
X_test = np.reshape(X_test, X_test.shape + (1,))

inputs = Input(shape=(timesteps, input_dim))
encoded = LSTM(latent_dim)(inputs)

decoded = RepeatVector(timesteps)(encoded)
decoded = LSTM(input_dim, return_sequences=True)(decoded)

sequence_autoencoder = Model(inputs, decoded)
encoder = Model(inputs, encoded)


encoder.compile(optimizer='rmsprop', loss='binary_crossentropy')
encoder.fit(X_train, X_train,
        nb_epoch=nb_epoch,
        batch_size=batch_size,
        validation_data=(X_test, X_test))
