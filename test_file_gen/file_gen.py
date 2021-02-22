import sys
import getopt
import json
from struct import *

def help():
	print("Usage: python file_gen.py -c <config_path> -o <output_path>")
	print("<config_path> -- the json config contains file sturct.")
	print("<output_path> -- the output path")

# convert python value to c bytes, according to the type
def value_to_bytes(typ, length, val):
    if typ == "short":
        return pack('h', val)
    elif typ == "unsigned short":
        return pack('H', val)
    elif typ == "int" or type == "int32":
        return pack('i', val)
    elif typ == "uint" or type == "uint32":
        return pack('I', val)
    elif typ == "int64":
        return pack('q', val)
    elif typ == "uint64":
        return pack('Q', val)
    elif typ == "char" and length == 1:
        return pack('c', val.encode('utf-8'))
    elif typ == "char" and length > 1:
         return val.encode('utf-8')

def main(argv):
    try:
        opts, args = getopt.getopt(argv, "hc:o:", ["cfg=","out=",])
    except getopt.GetoptError:
        help()
        sys.exit(2)

    confpath = ""
    outpath = ""

    for opt, arg in opts:
        if opt == '-h':
            help()
            sys.exit()
        elif opt in ("-c", "--cfg"):
            confpath = arg
        elif opt in ("-o", "--out"):
            outpath = arg

    with open(confpath) as f:
        conf = json.load(f)
        bts = value_to_bytes("char", 10, "A123456789")
        print(bts)

if __name__ == "__main__":
    main(sys.argv[1:])


