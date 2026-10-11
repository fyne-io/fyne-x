package widget

const traceShaderBody = `
uniform vec2 frame;
uniform vec4 bounds;
uniform float time;
uniform float inset;
uniform float radius;
uniform float stroke;
uniform float speed;
uniform float primaryR;
uniform float primaryG;
uniform float primaryB;
uniform float buttonR;
uniform float buttonG;
uniform float buttonB;

float roundedRectDistance(vec2 p, vec2 halfSize, float r) {
    vec2 q = abs(p) - halfSize + vec2(r);
    return length(max(q, 0.0)) + min(max(q.x, q.y), 0.0) - r;
}

float wrappedDistance(float a, float b) {
    float d = abs(a - b);
    return min(d, 1.0 - d);
}

float pathPosition(vec2 p, vec2 size, float r) {
    float left = inset;
    float topY = inset;
    float right = size.x - inset;
    float bottom = size.y - inset;
    float leftC = left + r;
    float rightC = right - r;
    float topC = topY + r;
    float bottomC = bottom - r;
    float top = max(0.0, rightC - leftC);
    float side = max(0.0, bottomC - topC);
    float arc = 1.57079632679 * r;
    float perimeter = max(1.0, top * 2.0 + side * 2.0 + arc * 4.0);
    float s = 0.0;

    if (p.y <= topC && p.x >= leftC && p.x <= rightC) {
        s = p.x - leftC;
    } else if (p.x >= rightC && p.y <= topC) {
        vec2 d = p - vec2(rightC, topC);
        float a = clamp(atan(d.x, -d.y), 0.0, 1.57079632679);
        s = top + a * r;
    } else if (p.x >= rightC && p.y >= topC && p.y <= bottomC) {
        s = top + arc + p.y - topC;
    } else if (p.x >= rightC && p.y >= bottomC) {
        vec2 d = p - vec2(rightC, bottomC);
        float a = clamp(atan(d.y, d.x), 0.0, 1.57079632679);
        s = top + arc + side + a * r;
    } else if (p.y >= bottomC && p.x >= leftC && p.x <= rightC) {
        s = top + arc + side + arc + rightC - p.x;
    } else if (p.x <= leftC && p.y >= bottomC) {
        vec2 d = p - vec2(leftC, bottomC);
        float a = clamp(atan(-d.x, d.y), 0.0, 1.57079632679);
        s = top + arc + side + arc + top + a * r;
    } else if (p.x <= leftC && p.y >= topC && p.y <= bottomC) {
        s = top + arc + side + arc + top + arc + bottomC - p.y;
    } else {
        vec2 d = p - vec2(leftC, topC);
        float a = clamp(atan(-d.y, -d.x), 0.0, 1.57079632679);
        s = top + arc + side + arc + top + arc + side + a * r;
    }

    return fract(s / perimeter);
}

void main() {
    vec2 size = vec2(bounds.z - bounds.x, bounds.w - bounds.y);
    vec2 local = vec2(gl_FragCoord.x - bounds.x, frame.y - gl_FragCoord.y - bounds.y);
    vec2 center = size * 0.5;
    float r = min(max(0.0, radius - inset), min(size.x, size.y) * 0.5 - inset - 1.0);
    vec2 halfSize = center - vec2(inset);

    float dist = roundedRectDistance(local - center, halfSize, r);
    float border = 1.0 - smoothstep(stroke, stroke + 0.85, abs(dist));
    float glow = 1.0 - smoothstep(stroke, stroke + 3.4, abs(dist));
    float t = pathPosition(local, size, r);
    float phase = fract(time * speed);
    float headA = 1.0 - smoothstep(0.00, 0.052, wrappedDistance(t, phase));
    float headB = 1.0 - smoothstep(0.00, 0.052, wrappedDistance(t, fract(phase + 0.5)));
    float trail = max(headA, headB);
    float pulse = 0.18 + trail * 0.82;
    float alpha = border * 0.78 + glow * trail * 0.18;

    if (alpha < 0.01) {
        discard;
    }

    vec3 primaryColor = vec3(primaryR, primaryG, primaryB);
    vec3 buttonColor = vec3(buttonR, buttonG, buttonB);
    vec3 color = mix(buttonColor, primaryColor, pulse);
    gl_FragColor = vec4(color, alpha);
}
`

const traceShaderSource = `#version 110
` + traceShaderBody

const traceShaderSourceES = `#version 100
#ifdef GL_ES
# ifdef GL_FRAGMENT_PRECISION_HIGH
precision highp float;
# else
precision mediump float;
# endif
precision mediump int;
#endif
` + traceShaderBody
