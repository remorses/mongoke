from vprof import runner
import time
from mongoke import main
# def hard():
#     time.sleep(2)
#     return 'ok'

# def lesshard():
#     time.sleep(1)

# def main(a, b):
#     print(hard())
#     lesshard()

runner.run(main, 'cmhp', args=('confs/testing.yaml', True), host='localhost', port=8000)