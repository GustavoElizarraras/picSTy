import tensorflow_hub as hub
import tensorflow as tf
from matplotlib import pyplot as plt
import sys

def load_image(path):
    img = tf.io.read_file(path)
    img = tf.image.decode_image(img, channels = 3)
    img = tf.image.convert_image_dtype(img, tf.float32)
    img = img[tf.newaxis, :]
    return img

if __name__ == "__name__":
    processed_path = '../go_web/processed_image'
    uploaded_image = sys.argv[2]
    selected_artwork = sys.argv[3]
    model = hub.load('https://tfhub.dev/google/magenta/arbitrary-image-stylization-v1-256/2')
    content = load_image(uploaded_image)
    style = load_image(selected_artwork)
    stylized_image = model(tf.constant(content), tf.constant(style))[0]
    tf.keras.utils.save_img(
        processed_path, stylized_image, data_format=None, file_format='.jpg', scale=True)

