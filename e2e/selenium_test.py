from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
import time


driver = webdriver.Edge()


BASE_URL = "http://localhost:8081"

try:
    #LOGIN
    print("Testing User Login...")
    driver.get(f"{BASE_URL}/login")
    time.sleep(2)


    driver.find_element(By.NAME, "email").send_keys("suhansun13@gmail.com")
    driver.find_element(By.NAME, "password").send_keys("123123")
    driver.find_element(By.ID, "login_button").click()

    driver.get(f"{BASE_URL}/posts/new")

    time.sleep(2)
    assert "Create New Post" in driver.page_source

    print("All tests passed!")

except Exception as e:
    print(f"Test failed: {e}")

finally:

    driver.quit()
