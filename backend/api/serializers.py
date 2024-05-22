from marshmallow import Schema, fields


class NodeMeasurementSchema(Schema):
    class Meta:
        ordered = True

    measurement_id = fields.Int()
    value = fields.Float()
    timestamp = fields.DateTime()
    node_id = fields.Int()
    type = fields.String()


class NodeSchema(Schema):
    node_id = fields.Int()
    mac_address = fields.String()
    name = fields.String()
    access_point_id = fields.Int()
    measurements = fields.Nested(NodeMeasurementSchema, many=True)

