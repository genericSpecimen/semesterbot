import telegram
from telegram.ext import Updater, CommandHandler
from telegram.utils.helpers import escape_markdown
import logging

from selenium import webdriver
import difflib
import os

def start(update, context):
    context.bot.send_message(chat_id=update.effective_chat.id,
                             text="Hello, I am a bot that will help with your semester needs.")

def status(update, context):
    update.message.reply_text(reply_to_message_id=update.message.message_id, text="I'm ok, thank for concern.")

def helptext(update, context):
    help_text = """Type one of the following commands:
    /help - display this message.
    /status - check if I'm ok.
    /check - check the website for news.
    /caps text - convert text to T E X T.
    """
    update.message.reply_text(reply_to_message_id=update.message.message_id, text=help_text)

def caps(update, context):
    text_caps = ' '.join(context.args).upper()
    spaced_text = ' '.join(text_caps)
    update.message.reply_text(reply_to_message_id=update.message.message_id, text=spaced_text)

def check_website(update, context):
    options = webdriver.FirefoxOptions()
    options.add_argument('--headless')
    driver = webdriver.Firefox(executable_path='./geckodriver', options=options)
    driver.get('http://cs.du.ac.in/')

    tables = driver.find_elements_by_tag_name('table')

    allnews = tables[0].find_elements_by_tag_name('td')
    allnews += tables[1].find_elements_by_tag_name('td')

    with open("new.md", "w") as newfile:
        print('# News', file=newfile)
        for news in allnews:
            news_text = escape_markdown(news.text, version=2)
            print(f'*{news_text}*', file=newfile)
            links = news.find_elements_by_tag_name('a')
            for link in links:
                link_text = escape_markdown(link.text, version=2)                
                link_href = escape_markdown(link.get_attribute("href"), version=2)
                print(f'[{link_text}]({link_href})', file=newfile)
            print(file=newfile)

    with open("new.md", "r") as newfile, open("old.md", "r") as oldfile, open("diff.md", "w") as difffile:
        old = oldfile.readlines()
        new = newfile.readlines()
        
        for line in difflib.unified_diff(old, new):
            if (line.startswith('+') and not line.startswith('+++')):
                print(line[1:], end="", file=difffile)
        
        os.rename("new.md", "old.md")

    with open("diff.md", "r") as difffile:
        diff = difffile.readlines()
        #print(diff)
        if (len(diff) > 0):
            text_message = ''.join(diff)
            print(''.join(diff))
            
        else:
            text_message = "Already up to date\."
            print("Already up to date\.")

    os.remove("diff.md")

    update.message.reply_markdown_v2(reply_to_message_id=update.message.message_id, text=text_message)

if __name__ == '__main__':
    TOKEN = os.getenv('TG_TOKEN')
    logging.basicConfig(format='%(asctime)s - %(name)s - %(levelname)s - %(message)s', level=logging.INFO)

    bot = telegram.Bot(token=TOKEN)
    #print(bot.get_me())

    updater = Updater(token=TOKEN, use_context=True)

    start_handler = CommandHandler('start', start)
    updater.dispatcher.add_handler(start_handler)

    status_handler = CommandHandler('status', status)
    updater.dispatcher.add_handler(status_handler)
    
    help_handler = CommandHandler('help', helptext)
    updater.dispatcher.add_handler(help_handler)

    caps_handler = CommandHandler('caps', caps)
    updater.dispatcher.add_handler(caps_handler)

    check_website_handler = CommandHandler('check', check_website)
    updater.dispatcher.add_handler(check_website_handler)

    updater.start_polling()
    updater.idle()
    
