# The Más Bandwidth Source License (MBSL)

The license for new Más Bandwidth libraries: **BSD 3-Clause plus exactly one
clause,** so products that incorporate the library must credit it in their
credits. Free to use, source open, credit required.

Scope, plainly:

- **New libraries** ship under the MBSL from their first release.
- **Ports with no external contributors** (netcode.go, netcode.rs,
  reliable.go, reliable.rs, serialize.go) carry it too.
- **The existing libraries keep their licenses unchanged:** yojimbo,
  netcode, reliable, serialize stay BSD 3-Clause and fixed3d stays MIT,
  forever. Roughly seventy people contributed to those under those terms,
  and that's honored. For them, crediting is an official request and the
  expected standard, not a license term.

## The license

```
Más Bandwidth Source License (MBSL)

Copyright © <years>, Más Bandwidth LLC

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice,
   this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Any product that incorporates this software must include the following
   credit in the product's credits, or in its documentation if the product
   has no credits, under the heading "Más Bandwidth LLC":

       <library> by Glenn Fiedler

   If this software depends on other Más Bandwidth libraries, the same
   credit must be included for each of them, as listed in this repository's
   README.

4. Neither the name of the copyright holder nor the names of its
   contributors may be used to endorse or promote products derived from
   this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
```

Clauses 1, 2 and 4 and the disclaimer are the BSD 3-Clause license, verbatim.
Clause 3 is the whole difference. Ports of BSD-licensed Más Bandwidth
libraries append a short notice retaining the upstream BSD terms for derived
portions.

*Known trade-offs, stated openly: the credit clause makes the MBSL
GPL-incompatible and outside the OSI-approved list. That's the deliberate
price of "keep people honest." The BSD/MIT libraries are unaffected.*
