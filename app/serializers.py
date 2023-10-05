from marshmallow import Schema, fields


class MeasurementSchema(Schema):
    id = fields.Int()
    temperature = fields.Float()
    humidity = fields.Float()
    timestamp = fields.DateTime()


class CompostBinSchema(Schema):
    id = fields.Int()
    name = fields.String()
    description = fields.String()
    measurements = fields.Nested(MeasurementSchema, many=True)

