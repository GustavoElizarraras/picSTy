#!/usr/bin/env python3

import tensorflow_hub as hub
import tensorflow as tf
import sys
import os
# os.environ["CUDA_VISIBLE_DEVICES"] = "-1"
# os.environ['TF_CPP_MIN_LOG_LEVEL'] = '3'
def load_image(path):
    img = tf.io.read_file(path)
    img = tf.image.decode_image(img, channels = 3)
    img = tf.image.convert_image_dtype(img, tf.float32)
    img = img[tf.newaxis, :]
    return img

def main():

    processed_path = '../go_web/processed_image' + sys.argv[1]
    uploaded_image = sys.argv[1]
    selected_artwork = sys.argv[2]
    model = hub.load('https://tfhub.dev/google/magenta/arbitrary-image-stylization-v1-256/2')
    content = load_image(uploaded_image)
    style = load_image(selected_artwork)
    stylized_image = model(tf.constant(content), tf.constant(style))[0]
    tf.keras.preprocessing.image.save_img(
        processed_path, tf.squeeze(stylized_image), data_format='channels_last', file_format='png', scale=True)


if __name__ == "__main__": 
    main()