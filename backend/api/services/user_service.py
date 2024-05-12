from ..models import User, AccessPoint


def get_user_by_id(user_id):
    try:
        user = User.query.get(user_id)
        return user
    except Exception as e:
        raise e


def validate_access_point_from_user(access_point_id, user_id):
    access_point = AccessPoint.query.get(access_point_id)
    if not access_point or access_point.user_id != user_id:
        raise Exception("invalid user")
