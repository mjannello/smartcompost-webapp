from marshmallow import Schema, fields


class MeasurementSchema(Schema):
    class Meta:
        ordered = True

    measurement_id = fields.Int()
    value = fields.Float()
    timestamp = fields.DateTime()
    compost_bin_id = fields.Int()
    type = fields.String()


class CompostBinSchema(Schema):
    compost_bin_id = fields.Int()
    name = fields.String()
    description = fields.String()
    measurements = fields.Nested(MeasurementSchema, many=True)

