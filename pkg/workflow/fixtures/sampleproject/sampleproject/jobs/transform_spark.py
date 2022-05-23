from awsglue.context import GlueContext
from pyspark import SparkContext

if __name__ == "__main__":
    glue = GlueContext(SparkContext.getOrCreate())
