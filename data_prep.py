import numpy as np
import pandas as pd
import tensorflow as tf
from tqdm import tqdm

import utils

tf.enable_eager_execution()

df = pd.read_csv('data/labels.csv')

train_writer = tf.python_io.TFRecordWriter('data/dogs224_train.tfrecord')
valid_writer = tf.python_io.TFRecordWriter('data/dogs224_valid.tfrecord')
ws = np.array([train_writer, valid_writer])
for _, row in tqdm(df.iterrows()):
    fname = "data/train/" + row.id + ".jpg"
    data = tf.read_file(fname)
    image = utils.preprocess_image(data, (224, 224))
    ex = utils.make_image_example(
        image.numpy().flatten(), row.breed)

    # Randomly choose whether the example should be training or validation
    w = np.random.choice(ws, p=[.8, .2])
    w.write(ex.SerializeToString())

train_writer.close()
valid_writer.close()

# Also save the label vocabulary
labels_vocab = df.breed.unique()
np.save('data/labelvocab', labels_vocab)
