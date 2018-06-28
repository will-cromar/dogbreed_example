import tensorflow as tf


def preprocess_image(image_bytes, size):
    im = tf.image.decode_image(image_bytes)
    im = tf.image.convert_image_dtype(im, tf.float32)
    im = tf.image.resize_bilinear(
          [im], size, align_corners=True)[0]

    return im


def make_image_example(flat_image, label):
    features = {}
    features["image"] = tf.train.Feature(
        float_list=tf.train.FloatList(value=flat_image))
    features["label"] = tf.train.Feature(
        bytes_list=(tf.train.BytesList(
            value=[label.encode('utf-8')])))

    ex = tf.train.Example(
        features=tf.train.Features(
            feature=features,
        )
    )

    return ex


def decode_image_example(ex, size):
    features = {}
    features["image"] = tf.FixedLenFeature(size, tf.float32)
    features["label"] = tf.FixedLenFeature((), tf.string)

    x = tf.parse_single_example(ex, features)
    y = x.pop("label")
    return x, y
