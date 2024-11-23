#version 460 core
in vec2 local_pos;
in float op_code;
in float radius;
in vec4 color;
in float width;
in float height;
in vec2 tex_coord;
in float texture_index;
in float font_index;

out vec4 fragColor;

uniform sampler2D samplers[24];

const float OP_CODE_VERTEX = 1.0;
const float OP_CODE_CIRCLE = 2.0;
const float OP_CODE_RECT = 3.0;
const float OP_CODE_TEXT = 4.0;
const float OP_CODE_TEXTURE = 5.0;

float sdCircle(vec2 p, float r) {
    return length(p) - r;
}

float sdRoundedRect(vec2 p, vec2 bounds, float r) {
    vec2 b = bounds - vec2(r);
    vec2 q = abs(p) - b;
    return length(max(q, 0.0)) - r;
}

float sdEquilateralTriangle(vec2 p) {
    float k = sqrt(3.0);
    p.x = abs(p.x) - 0.5;
    p.y = p.y + 0.5 / k;
    if (p.x + k * p.y > 0.0) p = vec2(p.x - k * p.y, -k * p.x - p.y) / 2.0;
    p.x -= clamp(p.x, -1.0, 0.0);
    return -length(p) * sign(p.y);
}

void main() {
    vec2 p = local_pos;
    fragColor = vec4(color.rgb, 0.0);

    if (op_code == OP_CODE_VERTEX) {
        fragColor = color;
    } else if (op_code == OP_CODE_CIRCLE) {
        float sdf = sdCircle(p, radius);
        if (sdf < 0.0) {
            fragColor = color;
        }
    } else if (op_code == OP_CODE_RECT) {
        float sdf = sdRoundedRect(p, vec2(width, height) * 0.5, radius);
        if (sdf < 0.0) {
            fragColor = color;
        }
    }
    if (op_code == OP_CODE_TEXT) {
        int idx = int(font_index);
        float alpha = texture(samplers[idx], tex_coord).a;
        fragColor = vec4(color.rgb, alpha * color.a);
    }

    if (op_code == OP_CODE_TEXTURE) {
        int idx = int(texture_index);
        fragColor = texture(samplers[idx], tex_coord);
    }
}
