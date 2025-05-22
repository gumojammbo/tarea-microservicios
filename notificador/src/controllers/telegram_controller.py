##!/usr/bin/env python
# -*- coding: utf-8 -*-
#-------------------------------------------------------------------------
# Archivo: telegram_controller.py
# Capitulo: Estilo Microservicios
# Autor(es): Perla Velasco & Yonathan Mtz. & Jorge Solís
# Version: 3.0.0 Febrero 2022
# Descripción:
#
#   Ésta clase define el controlador del microservicio API. 
#   Implementa la funcionalidad y lógica de negocio del Microservicio.
#
#   A continuación se describen los métodos que se implementaron en ésta clase:
#
#                                             Métodos:
#           +------------------------+--------------------------+-----------------------+
#           |         Nombre         |        Parámetros        |        Función        |
#           +------------------------+--------------------------+-----------------------+
#           |     send_message()     |         Ninguno          |  - Procesa el mensaje |
#           |                        |                          |    recibido en la     |
#           |                        |                          |    petición y ejecuta |
#           |                        |                          |    el envío a         |
#           |                        |                          |    Telegram.          |
#           +------------------------+--------------------------+-----------------------+
#
#-------------------------------------------------------------------------
from flask import request, jsonify
import json

from telegram import Bot
from telegram.error import TelegramError
import asyncio
from src.helpers.config import load_config

class TelegramController:

    @staticmethod
    def send_message():
        data = json.loads(request.data)
        if not data:
            return jsonify({"msg": "invalid request"}), 400
        
        # Cargar credenciales
        config = load_config()
        TELEGRAM_TOKEN = config['TELEGRAM']['TOKEN']
        CHAT_ID = config['TELEGRAM']['CHAT_ID']

        # Cargar datos del mensaje
        message = data["message"]

        
        async def enviar():
            # Crear bot
            bot = Bot(token=TELEGRAM_TOKEN)
            # Enviar mensaje
            return await bot.send_message(text=message, chat_id=CHAT_ID)

        # Crear bot y enviar mensaje
        try:
            msg = asyncio.run(enviar())
            print(msg.message_id, flush=True)
            return jsonify({"msg": "Message sent successfully"}), 200
        except TelegramError as e:
            print(e, flush=True)
            return jsonify({"msg": "Failed to send message"}), 500