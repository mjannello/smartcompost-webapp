from ..models import User


def get_user_by_id(user_id):
    try:
        user = User.query.get(user_id)
        return user
    except Exception as e:
        raise e
