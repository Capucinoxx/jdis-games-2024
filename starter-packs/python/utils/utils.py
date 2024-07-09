import uuid

def read_string_until_null(byte_array, end_index=None):
    # Trouver l'index du premier caractère nul
    if end_index is None:
        end_index = byte_array.find(b'\0')
        if end_index == -1:
            return None  # Retourner None si aucun caractère nul n'est trouvé
    
    # Extraire la chaîne de caractères jusqu'au caractère nul
    string = byte_array[:end_index].decode('utf-8')
    return string, end_index


def read_uuid(byte_array, end_index):
    return str(uuid.UUID(bytes=byte_array[:end_index]))