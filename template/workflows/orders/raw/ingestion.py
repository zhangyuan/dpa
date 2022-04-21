import sys,os
sys.path.append(os.getcwd())

import vendor.helper as helper
import pandas as pd

class Ingestion(object):
    def __init__(self) -> None:
        pass

    def extract(self) -> pd.DataFrame:
        print("extract data...")
        return pd.DataFrame()

    def transform(self, raw_df: pd.DataFrame) -> pd.DataFrame:
        print("extract df...")
        return pd.DataFrame()

    def load(self, df: pd.DataFrame):
        print("write df...")

    def perform(self):
        raw_df = self.extract()
        df = self.transform(raw_df)
        self.load(df)

if __name__ == "__main__":
    print(helper.hello())
    ingestion = Ingestion()
    ingestion.perform()
