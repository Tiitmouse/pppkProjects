import os
import time
import gzip
import shutil

from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.firefox.options import Options

download_dir = os.path.join(os.getcwd(), "dwnData")
if not os.path.exists(download_dir):
    os.makedirs(download_dir)

options = Options()
options.set_preference("browser.download.folderList", 2)
options.set_preference("browser.download.manager.showWhenStarting", False)
options.set_preference("browser.download.dir", download_dir)
options.set_preference("browser.helperApps.neverAsk.saveToDisk", "application/gzip")

driver = webdriver.Firefox(options=options)

def uncompress_gz_file(gz_file_path, extract_dir):
    with gzip.open(gz_file_path, 'rb') as f_in:
        with open(os.path.join(extract_dir, os.path.basename(gz_file_path).replace('.gz', '') + '.txt'), 'wb') as f_out:
            shutil.copyfileobj(f_in, f_out)

try:
    driver.get("https://xenabrowser.net/datapages/?hub=https://tcga.xenahubs.net:443")
    time.sleep(5)
    cohorts = [element.get_attribute("href") for element in
               driver.find_elements(By.CSS_SELECTOR, "li.MuiTypography-root-158 > a")]

    for cohort in cohorts:
        driver.get(cohort)
        time.sleep(2)

        hubs = driver.find_elements(By.XPATH, "//a[contains(text(), 'pancan normalized')]")
        for hub in hubs:
            details_page = hub.get_attribute("href")
            driver.get(details_page)
            time.sleep(5)
            try:
                download_link = driver.find_element(By.XPATH, "/html/body/div/div[2]/div/div/div/span[6]/span/a[1]")
                download_link.click()
                while any(fname.endswith('.part') for fname in os.listdir(download_dir)):
                    time.sleep(1)
                for fname in os.listdir(download_dir):
                    if fname.endswith('.gz'):
                        uncompress_gz_file(os.path.join(download_dir, fname), download_dir)
            except:
                print("Download link not found")

        driver.back()
        time.sleep(2)

    while any(fname.endswith('.part') for fname in os.listdir(download_dir)):
        time.sleep(1)

finally:
    driver.quit()