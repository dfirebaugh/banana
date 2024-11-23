#version 460 core

layout(location = 0) in vec2 in_pos;
layout(location = 1) in vec2 in_shape_pos;
layout(location = 2) in vec2 in_local_pos;
layout(location = 3) in float in_op_code;
layout(location = 4) in float in_radius;
layout(location = 5) in vec4 in_color;
layout(location = 6) in float in_width;
layout(location = 7) in float in_height;
layout(location = 8) in vec2 in_tex_coord;
layout(location = 9) in vec2 in_resolution;
layout(location = 10) in float in_texture_index;
layout(location = 11) in float in_font_index;

out vec2 local_pos;
out float op_code;
out float radius;
out vec4 color;
out float width;
out float height;
out vec2 tex_coord;
out uint stencil_write_value;
out uint stencil_test_value;
out float texture_index;
out float font_index;

void main() {
    vec2 scaledShapePos = in_shape_pos;
    if (in_op_code == 4.0 || in_op_code == 5.0) {
        gl_Position = vec4(in_pos, 0.0, 1.0);
    } else {
        gl_Position = vec4(scaledShapePos + in_local_pos / in_resolution * 2.0, 0.0, 1.0);
    }
    local_pos = in_local_pos;
    op_code = in_op_code;
    radius = in_radius;
    color = in_color;
    width = in_width;
    height = in_height;
    tex_coord = in_tex_coord;
    texture_index = in_texture_index;
    font_index = in_font_index;
}
