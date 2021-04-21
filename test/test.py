import os
import unittest

class TestNetcat(unittest.TestCase):
    def test_message_ok(self):
        os.system("docker build -f ./Dockerfile -t \"test:latest\" .")
        process = os.popen("docker run --network 7574-tp0_testing_net test:latest")
        
        message = process.readline().rstrip()
        self.assertEqual(message, "Your Message has been received: b\'message\'")

        process.close()
